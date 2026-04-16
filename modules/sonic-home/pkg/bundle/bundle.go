package bundle

import (
	"archive/tar"
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/sonic-family/sonic-home/pkg/manifest"
)

const (
	headerPath    = "header.sheh"
	manifestPath  = "manifest.json"
	signaturePath = "signature.sig"
)

var requiredEntries = []string{
	headerPath,
	manifestPath,
	signaturePath,
	"payload/",
	"payload/scripts/",
	"payload/scripts/pre-install.sh",
	"payload/scripts/post-install.sh",
	"payload/scripts/healthcheck.sh",
}

func BuildDraftBundle(outPath string, m manifest.Manifest) error {
	if strings.TrimSpace(outPath) == "" {
		return fmt.Errorf("output path is required")
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	tw := tar.NewWriter(f)
	defer tw.Close()

	header := []byte("SHE-DRAFT-1\n")
	signature := []byte("unsigned\n")

	if err := writeFile(tw, headerPath, header, 0o644); err != nil {
		return err
	}

	if err := writeDir(tw, "payload/", 0o755); err != nil {
		return err
	}
	if err := writeDir(tw, "payload/base/", 0o755); err != nil {
		return err
	}
	if err := writeDir(tw, "payload/scripts/", 0o755); err != nil {
		return err
	}
	if err := writeFile(tw, "payload/scripts/pre-install.sh", []byte("#!/bin/sh\nset -e\n"), 0o755); err != nil {
		return err
	}
	if err := writeFile(tw, "payload/scripts/post-install.sh", []byte("#!/bin/sh\nset -e\n"), 0o755); err != nil {
		return err
	}
	if err := writeFile(tw, "payload/scripts/healthcheck.sh", []byte("#!/bin/sh\nexit 0\n"), 0o755); err != nil {
		return err
	}

	components, err := writeSourcePayload(tw, m.SourcePath, outPath)
	if err != nil {
		return err
	}
	m.Components = components

	manifestBytes, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return err
	}
	manifestBytes = append(manifestBytes, '\n')

	if err := writeFile(tw, manifestPath, manifestBytes, 0o644); err != nil {
		return err
	}
	if err := writeFile(tw, signaturePath, signature, 0o644); err != nil {
		return err
	}

	return nil
}

func VerifyDraftBundle(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	tr := tar.NewReader(f)
	found := make(map[string]bool)
	var manifestBody bytes.Buffer
	checksums := make(map[string]string)
	sizes := make(map[string]int64)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		name := hdr.Name
		if hdr.FileInfo().IsDir() && !strings.HasSuffix(name, "/") {
			name += "/"
		}
		found[name] = true

		if name == manifestPath {
			if _, err := io.Copy(&manifestBody, tr); err != nil {
				return err
			}
			continue
		}

		if hdr.FileInfo().IsDir() {
			continue
		}

		h := sha256.New()
		n, err := io.Copy(h, tr)
		if err != nil {
			return err
		}
		checksums[name] = "sha256:" + hex.EncodeToString(h.Sum(nil))
		sizes[name] = n
	}

	for _, entry := range requiredEntries {
		if !found[entry] {
			return fmt.Errorf("missing required entry: %s", entry)
		}
	}

	var parsed manifest.Manifest
	if err := json.Unmarshal(manifestBody.Bytes(), &parsed); err != nil {
		return fmt.Errorf("invalid manifest.json: %w", err)
	}
	if parsed.SchemaVersion == "" || parsed.BundleID == "" {
		return fmt.Errorf("manifest missing required identifiers")
	}
	for _, c := range parsed.Components {
		ck, ok := checksums[c.Source]
		if !ok {
			return fmt.Errorf("component source missing: %s", c.Source)
		}
		if c.Checksum != "" && c.Checksum != ck {
			return fmt.Errorf("component checksum mismatch for %s", c.Source)
		}
		if c.SizeBytes > 0 && c.SizeBytes != sizes[c.Source] {
			return fmt.Errorf("component size mismatch for %s", c.Source)
		}
	}

	return nil
}

func writeSourcePayload(tw *tar.Writer, sourcePath string, outPath string) ([]manifest.Component, error) {
	root := strings.TrimSpace(sourcePath)
	if root == "" {
		root = "."
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	absOut, err := filepath.Abs(outPath)
	if err != nil {
		return nil, err
	}

	type entry struct {
		rel  string
		full string
		size int64
	}
	var entries []entry
	if err := filepath.WalkDir(absRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == ".cursor" {
				return filepath.SkipDir
			}
			return nil
		}
		if path == absOut {
			return nil
		}
		if !d.Type().IsRegular() {
			return nil
		}

		rel, err := filepath.Rel(absRoot, path)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		entries = append(entries, entry{rel: filepath.ToSlash(rel), full: path, size: info.Size()})
		return nil
	}); err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool { return entries[i].rel < entries[j].rel })

	components := make([]manifest.Component, 0, len(entries))
	for _, e := range entries {
		payloadPath := "payload/base/" + e.rel
		sum, err := writeFileFromDisk(tw, payloadPath, e.full, 0o644)
		if err != nil {
			return nil, err
		}
		components = append(components, manifest.Component{
			Name:      e.rel,
			Type:      "file",
			Source:    payloadPath,
			Checksum:  "sha256:" + sum,
			SizeBytes: e.size,
		})
	}

	return components, nil
}

func writeDir(tw *tar.Writer, path string, mode int64) error {
	hdr := &tar.Header{
		Name:     path,
		Typeflag: tar.TypeDir,
		Mode:     mode,
	}
	return tw.WriteHeader(hdr)
}

func writeFile(tw *tar.Writer, path string, body []byte, mode int64) error {
	hdr := &tar.Header{
		Name: path,
		Mode: mode,
		Size: int64(len(body)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	_, err := tw.Write(body)
	return err
}

func writeFileFromDisk(tw *tar.Writer, targetPath string, diskPath string, mode int64) (string, error) {
	f, err := os.Open(diskPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return "", err
	}

	hdr := &tar.Header{
		Name: targetPath,
		Mode: mode,
		Size: info.Size(),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return "", err
	}

	h := sha256.New()
	w := io.MultiWriter(tw, h)
	if _, err := io.Copy(w, f); err != nil {
		return "", err
	}
	return hashHex(h), nil
}

func hashHex(h hash.Hash) string {
	return hex.EncodeToString(h.Sum(nil))
}

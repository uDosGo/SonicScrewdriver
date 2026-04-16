# bundle package

Draft `.she` archive writer and verifier.

Current capabilities:
- create a draft tar-based `.she` layout with:
  - `header.sheh`
  - `manifest.json`
  - `signature.sig` (placeholder)
  - `payload/scripts/*` placeholders
- verify required entries and parse `manifest.json`

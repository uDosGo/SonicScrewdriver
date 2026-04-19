package library

// ManifestSchema defines the expected structure of a game manifest
const ManifestSchema = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "minLength": 1,
      "description": "The name of the game"
    },
    "version": {
      "type": "string",
      "pattern": "^\\d+\\.\\d+\\.\\d+(-[a-zA-Z0-9]+)?$",
      "description": "Semantic version of the game"
    },
    "title": {
      "type": "string",
      "description": "Display title of the game"
    },
    "description": {
      "type": "string",
      "description": "Description of the game"
    },
    "developer": {
      "type": "string",
      "description": "Game developer"
    },
    "publisher": {
      "type": "string",
      "description": "Game publisher"
    },
    "release_date": {
      "type": "string",
      "description": "Release date"
    },
    "genres": {
      "type": "array",
      "items": {
        "type": "string"
      },
      "description": "Game genres"
    },
    "platforms": {
      "type": "array",
      "items": {
        "type": "string",
        "enum": ["windows", "macos", "linux", "docker"]
      },
      "description": "Supported platforms"
    },
    "container": {
      "type": "object",
      "properties": {
        "image": {
          "type": "string",
          "description": "Docker image name"
        },
        "tag": {
          "type": "string",
          "description": "Docker image tag"
        },
        "ports": {
          "type": "array",
          "items": {
            "type": "string",
            "pattern": "^\\d+(:\\d+)?$"
          },
          "description": "Port mappings"
        },
        "volumes": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Volume mounts"
        },
        "env": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          },
          "description": "Environment variables"
        }
      },
      "required": ["image"]
    },
    "requirements": {
      "type": "object",
      "properties": {
        "min_cpu": {
          "type": "string",
          "description": "Minimum CPU requirement"
        },
        "min_memory": {
          "type": "string",
          "description": "Minimum memory requirement"
        },
        "min_storage": {
          "type": "string",
          "description": "Minimum storage requirement"
        },
        "gpu": {
          "type": "boolean",
          "description": "GPU required"
        }
      }
    },
    "assets": {
      "type": "object",
      "properties": {
        "icon": {
          "type": "string",
          "description": "Path to icon file"
        },
        "screenshots": {
          "type": "array",
          "items": {
            "type": "string"
          },
          "description": "Paths to screenshot files"
        }
      }
    }
  },
  "required": ["name", "version", "container"]
}`
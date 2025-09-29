#!/bin/bash
set -e

VERSION_FILE="internal/version/version.go"

if [ $# -ne 1 ]; then
  echo "Usage: $0 [patch|minor|major]"
  exit 1
fi

BUMP_TYPE="$1"
if [[ "$BUMP_TYPE" != "patch" && "$BUMP_TYPE" != "minor" && "$BUMP_TYPE" != "major" ]]; then
  echo "Invalid argument: $BUMP_TYPE"
  echo "Usage: $0 [patch|minor|major]"
  exit 1
fi

if [ ! -f "$VERSION_FILE" ]; then
  echo "Version file not found: $VERSION_FILE"
  exit 1
fi

# Extract current version from the file
CURRENT_VERSION=$(grep -Eo 'const Version = "[0-9]+\.[0-9]+\.[0-9]+"' "$VERSION_FILE" | grep -Eo '[0-9]+\.[0-9]+\.[0-9]+')
if [ -z "$CURRENT_VERSION" ]; then
  echo "Could not find current version in $VERSION_FILE"
  exit 1
fi

IFS='.' read -r MAJOR MINOR PATCH <<< "$CURRENT_VERSION"

case "$BUMP_TYPE" in
  patch)
    PATCH=$((PATCH + 1))
    ;;
  minor)
    MINOR=$((MINOR + 1))
    PATCH=0
    ;;
  major)
    MAJOR=$((MAJOR + 1))
    MINOR=0
    PATCH=0
    ;;
esac

NEW_VERSION="${MAJOR}.${MINOR}.${PATCH}"
NEW_TAG="v${NEW_VERSION}"

# Check for clean working directory
if [ -n "$(git status --porcelain)" ]; then
  echo "Working directory is not clean. Please commit or stash your changes first."
  exit 1
fi

# Update the version constant (cross-platform sed)
if [[ "$OSTYPE" == "darwin"* ]]; then
  sed -i '' -E "s/const Version = \".*\"/const Version = \"${NEW_VERSION}\"/" "$VERSION_FILE"
else
  sed -i -E "s/const Version = \".*\"/const Version = \"${NEW_VERSION}\"/" "$VERSION_FILE"
fi

echo "Updated $VERSION_FILE to version $NEW_VERSION"

# Commit the change
git add "$VERSION_FILE"
git commit -m "Chore: Bump version to $NEW_VERSION"

# Create tag
git tag "$NEW_TAG"
echo "Created tag $NEW_TAG"

# Optionally push (uncomment if you want to push automatically)
git push
git push origin "$NEW_TAG"

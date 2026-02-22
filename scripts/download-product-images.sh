#!/bin/bash

# Download product placeholder images from Wikimedia Commons.
# Source filter: CC0 only (public domain dedication).

set -euo pipefail

OUT_DIR="frontend/public/images/products"
MANIFEST="${OUT_DIR}/sources.csv"

mkdir -p "${OUT_DIR}"

cat > "${MANIFEST}" <<'CSV'
slug,file_name,source_image_url,source_page_url,license,license_url,artist
CSV

download_one() {
  local slug="$1"
  local query="$2"
  local token="$3"

  local raw_search="${query} incategory:\"CC-Zero\" filetype:bitmap"
  local encoded_search
  encoded_search="$(printf '%s' "${raw_search}" | jq -sRr @uri)"

  local offset=$((RANDOM % 200))
  local api_url="https://commons.wikimedia.org/w/api.php?action=query&format=json&generator=search&gsrsearch=${encoded_search}&gsrnamespace=6&gsrlimit=30&gsroffset=${offset}&prop=imageinfo&iiprop=url|extmetadata&iiurlwidth=1400"

  local json
  json="$(curl -fsSL "${api_url}")"

  if [[ "$(printf '%s' "${json}" | jq '.query.pages | length // 0')" -eq 0 ]]; then
    api_url="https://commons.wikimedia.org/w/api.php?action=query&format=json&generator=search&gsrsearch=${encoded_search}&gsrnamespace=6&gsrlimit=30&prop=imageinfo&iiprop=url|extmetadata&iiurlwidth=1400"
    json="$(curl -fsSL "${api_url}")"
  fi

  local entries entry_count idx selected_entry
  entries="$(
    printf '%s' "${json}" | jq -r --arg token "${token}" '
      (
        .query.pages
        | to_entries
        | map(select($token == "" or (.value.title | ascii_downcase | contains($token))))
      ) as $filtered
      | if ($filtered | length) > 0 then $filtered else (.query.pages | to_entries) end
      | .[]
      | @base64
    '
  )"

  entry_count="$(printf '%s\n' "${entries}" | sed '/^$/d' | wc -l | tr -d ' ')"
  if [[ "${entry_count}" -eq 0 ]]; then
    echo "Skipping ${slug}: no image candidates found" >&2
    return
  fi

  idx=$((RANDOM % entry_count + 1))
  selected_entry="$(printf '%s\n' "${entries}" | sed -n "${idx}p")"

  local image_url source_page license license_url artist file_name ext
  image_url="$(printf '%s' "${selected_entry}" | base64 --decode | jq -r '.value.imageinfo[0].thumburl // .value.imageinfo[0].url')"
  source_page="$(printf '%s' "${selected_entry}" | base64 --decode | jq -r '.value.imageinfo[0].descriptionurl')"
  license="$(printf '%s' "${selected_entry}" | base64 --decode | jq -r '.value.imageinfo[0].extmetadata.LicenseShortName.value // "Unknown"')"
  license_url="$(printf '%s' "${selected_entry}" | base64 --decode | jq -r '.value.imageinfo[0].extmetadata.LicenseUrl.value // ""')"
  artist="$(printf '%s' "${selected_entry}" | base64 --decode | jq -r '.value.imageinfo[0].extmetadata.Artist.value // ""' | sed 's/<[^>]*>//g' | tr -d '\r\n')"

  if [[ -z "${image_url}" || "${image_url}" == "null" ]]; then
    echo "Skipping ${slug}: no image URL found" >&2
    return
  fi

  ext="$(printf '%s' "${image_url}" | sed -E 's/.*\.([A-Za-z0-9]+)(\?.*)?$/\1/' | tr '[:upper:]' '[:lower:]')"
  case "${ext}" in
    jpg|jpeg|png|webp|gif|tif|tiff)
      ;;
    *)
      ext="jpg"
      ;;
  esac

  file_name="${slug}.${ext}"
  curl -fsSL "${image_url}" -o "${OUT_DIR}/${file_name}"

  printf '"%s","%s","%s","%s","%s","%s","%s"\n' \
    "${slug}" \
    "${file_name}" \
    "${image_url}" \
    "${source_page}" \
    "${license}" \
    "${license_url}" \
    "${artist}" >> "${MANIFEST}"
}

download_one "wireless-headphones" "wireless headphones" "headphone"
download_one "smart-watch" "smart watch" "watch"
download_one "laptop-backpack" "laptop backpack" "backpack"
download_one "cotton-t-shirt" "cotton t-shirt" "shirt"
download_one "denim-jeans" "denim jeans" "jean"
download_one "winter-jacket" "winter jacket" "jacket"
download_one "the-great-gatsby" "book cover" "book"
download_one "to-kill-a-mockingbird" "novel book" "book"
download_one "nineteen-eighty-four" "classic book" "book"

echo "Downloaded images to ${OUT_DIR}"
echo "Saved source/license manifest to ${MANIFEST}"

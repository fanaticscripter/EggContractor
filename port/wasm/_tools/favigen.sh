#!/usr/bin/env zsh
setopt errexit nounset pipefail
bordered=0
if [[ $1 == -b || $1 == --border ]]; then
    bordered=1
    shift
fi
color=$1
output_dir=$2
[[ -n $color && -n $output_dir ]] || {
    cat >&2 <<EOF
$0 [-b|--border] <color> <output_directory>

Example:

  $0 '#ff0000' favicon.ico

to generate a red version of the favicon.

-b, --border adds a gray border to the favicon (used for primary colors with too little contract to a light/dark background.)
EOF
    exit 1
}
tmpdir="$(mktemp -d)"
echo "working directory: $tmpdir"
output_svg=$output_dir/favicon.svg
if (( bordered )); then
    cat >$output_svg <<EOF
<svg xmlns="http://www.w3.org/2000/svg" viewBox="-96 -32 576 576"><path fill="$color" d="M192 0C86 0 0 214 0 320s86 192 192 192 192-86 192-192S298 0 192 0z" /><path fill="#6b7280" d="M192 0C86 0 0 214 0 320s86 192 192 192 192-86 192-192S298 0 192 0zm0 480c-88.22 0-160-71.78-160-160 0-100.72 82.24-288 160-288s160 187.28 160 288c0 88.22-71.78 160-160 160z" /></svg>
EOF
else
    cat >$output_svg <<EOF
<svg xmlns="http://www.w3.org/2000/svg" viewBox="-96 -32 576 576"><path fill="$color" d="M192 0C86 0 0 214 0 320s86 192 192 192 192-86 192-192S298 0 192 0z" /></svg>
EOF
fi
for size in 16 32 48; do
    convert -background transparent -density 300 $output_svg -resize ${size}x${size} $tmpdir/favicon-$size.png
done
optipng $tmpdir/favicon-*.png
output_ico=$output_dir/favicon.ico
convert $tmpdir/favicon-*.png $output_ico
output_png=$output_dir/favicon-180.png
convert -background transparent -density 300 $output_svg -resize 180x180 $output_png
echo "generated $output_svg"
echo "generated $output_ico"
echo "generated $output_png"

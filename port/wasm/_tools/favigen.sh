#!/usr/bin/env zsh
setopt errexit nounset pipefail
color=$1
output_dir=$2
[[ -n $color && -n $output_dir ]] || {
    cat >&2 <<EOF
$0 <color> <output_directory>

Example:

  $0 '#ff0000' favicon.ico

to generate a red version of the favicon.
EOF
    exit 1
}
here=$0:A:h
orig_svg=$here/../../../static/favicon.svg
orig_svg=$orig_svg:A
tmpdir="$(mktemp -d)"
echo "working directory: $tmpdir"
output_svg=$output_dir/favicon.svg
sed "s/fill=\"#6b7280\"/fill=\"$color\"/" <$orig_svg >$output_svg
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

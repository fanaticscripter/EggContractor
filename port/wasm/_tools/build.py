#!/usr/bin/env python3

import argparse
import hashlib
import json
import os
import pathlib
import shutil
import subprocess


srcdir = pathlib.Path("src")
distdir = pathlib.Path("dist")

pagename = pathlib.Path(os.getcwd()).name
parent_distdir = pathlib.Path("../dist")
finaldest = parent_distdir / pagename


def wasm_handler(args):
    env = os.environ.copy()
    env["GOOS"] = "js"
    env["GOARCH"] = "wasm"
    print("go build -o dist/app.wasm")
    subprocess.check_call(["go", "build", "-o", "dist/app.wasm"], env=env)
    path_with_hash_suffix = add_hash_suffix(pathlib.Path("dist/app.wasm"))
    with open("dist/manifest.wasm.json", "w") as fp:
        fp.write(json.dumps({"app.wasm": path_with_hash_suffix.name}, indent=2))


def html_handler(args):
    manifest = read_manifests()
    with srcdir.joinpath("index.html").open() as fp:
        content = fp.read()
    for src, dst in manifest.items():
        content = content.replace(f"@@{src}@@", dst)
    with distdir.joinpath("index.html").open("w") as fp:
        fp.write(content)


def dist_handler(args):
    if finaldest.exists():
        print(f"rm -rf {finaldest}")
        shutil.rmtree(finaldest)
    print(f"mkdir -p {finaldest}")
    finaldest.mkdir(parents=True, exist_ok=True)
    copylist = ["index.html", "wasm_exec.js"]
    copylist.extend(read_manifests().values())
    copylist = [distdir / f for f in copylist]
    print(f"cp {' '.join(str(f) for f in copylist)} {finaldest}")
    for f in copylist:
        shutil.copy(f, finaldest)


def add_hash_suffix(path):
    hasher = hashlib.sha256()
    with path.open("rb") as fp:
        while True:
            chunk = fp.read(65536)
            if not chunk:
                break
            hasher.update(chunk)
    hash = hasher.hexdigest()
    path_with_hash_suffix = path.with_suffix(f".{hash[:8]}{path.suffix}")
    print(f"mv {path} {path_with_hash_suffix}")
    path.rename(path_with_hash_suffix)
    return path_with_hash_suffix


def read_manifests():
    manifest = dict()
    for f in distdir.glob("manifest*.json"):
        with f.open() as fp:
            manifest.update(json.load(fp))
    return manifest


def main():
    parser = argparse.ArgumentParser()
    subparsers = parser.add_subparsers(dest="cmd")
    subparsers.required = True

    wasm_parser = subparsers.add_parser("wasm")
    wasm_parser.set_defaults(handler=wasm_handler)

    html_parser = subparsers.add_parser("html")
    html_parser.set_defaults(handler=html_handler)

    dist_parser = subparsers.add_parser("dist")
    dist_parser.set_defaults(handler=dist_handler)

    args = parser.parse_args()
    args.handler(args)


if __name__ == "__main__":
    main()

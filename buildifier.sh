#!/bin/sh

bazel run @io_bazel_buildifier//buildifier -- \
    $(find $HOME/dorkalonius -iname BUILD -type f)
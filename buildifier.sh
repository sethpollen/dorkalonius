#!/bin/sh

bazel run @io_bazel_buildifier//buildifier -- \
    $(find $HOME/sbp/sbp_linux_config -iname BUILD -type f)
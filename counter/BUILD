load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "word_stream.go",
    ],
)

go_test(
    name = "word_stream_test",
    srcs = ["word_stream_test.go"],
    deps = [":go_default_library"],
)
load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = [
        "ebook_reader.go",
    ],
)

go_test(
    name = "ebook_reader_test",
    srcs = ["ebook_reader_test.go"],
    deps = [":go_default_library"],
)
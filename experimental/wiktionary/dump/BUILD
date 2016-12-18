load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "page_extractor.go",
    ],
)

go_binary(
    name = "page_extractor_main",
    srcs = ["page_extractor_main.go"],
    deps = [":go_default_library"],
)

go_binary(
    name = "template_extractor_main",
    srcs = ["template_extractor_main.go"],
    deps = [":go_default_library"],
)

load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "go_binary")

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

go_binary(
    name = "counter_main",
    srcs = ["counter_main.go"],
    data = ["//wiktionary:inflections_xml_bz2"],
    deps = [
        ":go_default_library",
        "//util:go_default_library",
        "//gutenberg:go_default_library",
        "//wiktionary:go_default_library",
    ],
)

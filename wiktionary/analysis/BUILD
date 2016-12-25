load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

filegroup(
    name = "lua_scripts",
    srcs = glob(["*.lua"]),
)

go_library(
    name = "go_default_library",
    srcs = [
        "inflection.go",
        "page_extractor.go",
    ],
    data = [":lua_scripts"],
)

go_test(
    name = "inflection_test",
    srcs = ["inflection_test.go"],
    deps = [":go_default_library"],
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

go_binary(
    name = "inflection_main",
    srcs = ["inflection_main.go"],
    data = ["//wiktionary/analysis/input_data"],
    deps = [
        ":go_default_library",
        "//wiktionary:go_default_library",
    ],
)

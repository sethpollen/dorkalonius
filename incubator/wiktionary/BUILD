load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

filegroup(
    name = "lua_scripts",
    srcs = glob(["*.lua"]),
)

go_library(
    name = "go_default_library",
    srcs = [
        "inflection.go",
    ],
    data = [":lua_scripts"],
)

go_test(
    name = "inflection_test",
    srcs = ["inflection_test.go"],
    deps = [":go_default_library"],
)

go_binary(
    name = "inflection_main",
    srcs = ["inflection_main.go"],
    data = ["//incubator/wiktionary/dump/data"],
    deps = [":go_default_library"],
)

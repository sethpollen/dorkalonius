load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_binary", "go_test")

go_binary(
    name = "embed_inflections_main",
    srcs = ["embed_inflections_main.go"],
    deps = ["//incubator/wiktionary:go_default_library"],
)

genrule(
    name = "embed_inflections",
    srcs = ["inflections.xml"],
    outs = ["embedded_inflections.go"],
    cmd = "./$(location embed_inflections_main)" +
          "  --source_file=\"$(location inflections.xml)\"" +
          "  --dest_file=\"$@\"",
    tools = [":embed_inflections_main"],
)

go_library(
    name = "go_default_library",
    srcs = [":embed_inflections"],
    visibility = ["//visibility:public"],
    deps = ["//incubator/wiktionary:go_default_library"],
)

go_test(
    name = "embed_test",
    srcs = ["embed_test.go"],
    deps = [":go_default_library"],
)
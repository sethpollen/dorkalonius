load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library", "go_test")

go_library(
    name = "go_default_library",
    srcs = ["go_embed_encoder.go"],
)

go_binary(
    name = "go_embed_main",
    srcs = ["go_embed_main.go"],
    visibility = ["//visibility:public"],
    deps = [":go_default_library"],
)

go_binary(
    name = "go_embed_encoder_test_main",
    srcs = ["go_embed_encoder_test_main.go"],
    deps = [
        ":go_default_library",
    ],
)

genrule(
    name = "go_embed_encoder_test_genrule",
    outs = ["go_embed_encoder_test.go"],
    cmd = "./$(location go_embed_encoder_test_main)" +
          "  --output_file=\"$@\"",
    tools = [":go_embed_encoder_test_main"],
)

go_test(
    name = "go_embed_encoder_test",
    srcs = [":go_embed_encoder_test_genrule"],
)

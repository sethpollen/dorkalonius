# Bazel rule for embedding data into a Go program.

# Arguments:
#   name: Name of the rule.
#   data: List of targets providing the files to embed.
#   package: Go package to put the generated code in.
def go_embed_data(name, data, package):
  native.genrule(
    name = name,
    srcs = data,
    outs = [name + "_embed_data.go"],
    cmd = "./$(location //tools:go_embed_main)" +
          "  --output_file=\"$@\"" +
          "  --package=\"" + package + "\"" +
          "  --method=\"Get_" + name + "\"" +
          "  $(SRCS)",
    tools = ["//tools:go_embed_main"],
  )
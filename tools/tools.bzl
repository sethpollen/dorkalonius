# Bazel rule for embedding data into a Go program.
#
# Arguments:
#   name: Name of the rule.
#   data: List of targets providing the files to embed.
#   package: Go package to put the generated code in.
def go_embed_data(
  name,
  data,
  package
):
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
  
# Bazel rule for converting a CSV file of word count data into a serialized
# WordSet object.
#
# Arguments:
#   name: Name of the rule
#   srcs: List of targets providing input CSV files
def word_set(
  name,
  srcs,
  csv_header_lines=0,
  csv_word_column=0,
  csv_weight_column=1
):
  native.genrule(
    name = name,
    srcs = srcs,
    outs = [name + ".wordset"],
    cmd = "./$(location //tools:csv_to_word_set_main)" +
          "  --output_file=\"$@\"" +
          "  --csv_header_lines=" + str(csv_header_lines) +
          "  --csv_word_column=" + str(csv_word_column) +
          "  --csv_weight_column=" + str(csv_weight_column) +
          "  $(SRCS)",
    tools = ["//tools:csv_to_word_set_main"],
  )

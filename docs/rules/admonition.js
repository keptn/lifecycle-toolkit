// @ts-check

"use strict";

module.exports = {
  "names": [ "github-admonition" ],
  "description": "asdfg",
  "information": new URL(
    "https://github.com/aepfli/markdownlint-rule-max-one-sentence-per-line"
  ),
  "tags": [ "admonition" ],
  "function": (params, onError) => {
    const admonitions = params.config.admonitions || [ "Note", "Warning" ];
    const admonitionStart = params.config.start || "> **";
    const admonitionEnd = params.config.end ||"**";

    const atBeginningRegex =
      new RegExp("^([\\n\\t ]*)(.*)(" +
        admonitions.join("|") +
        ")(:?\\*\\*)", "im");

    const atBeginning =
      (relevantToken) => {
        const check = atBeginningRegex.exec(relevantToken.line);
        if (check && (check[2] !== "> **" ||
          !admonitions.includes(check[3]) ||
          check[4] !== "**")) {
          const lengthForReplacement =
            check[2].length + check[3].length + check[4].length;
          const position = check[1].length + 1;
          onError({
            "lineNumber": relevantToken.lineNumber,
            "detail": relevantToken.line.substr(0, 10),
            "fixInfo": {
              "lineNumber": relevantToken.lineNumber,
              "editColumn": position,
              "deleteCount": lengthForReplacement,
              // eslint-disable-next-line max-len
              "insertText": admonitionStart + check[3].charAt(0).toUpperCase() + check[3].slice(1).toLowerCase() + admonitionEnd
            }
          });
        }
      };

    const wholeLineRegex =
      new RegExp("^([\\n\\t ]*)(.*)(" +
        admonitions.join("|") +
        "):?(.*)(\\*\\*)", "im");

    const wholeLine =
      (relevantToken) => {
        const check = wholeLineRegex
          .exec(relevantToken.line);
        if (check && (check[2] !== "> **" ||
          !admonitions.includes(check[3]) ||
          check[5] !== "**")) {
          const position = check[1] ? check[1].length + 1 : 1;

          const lengthForReplacement = relevantToken.line.length - position + 1;
          const replacement = admonitionStart +
            check[3].charAt(0).toUpperCase() + check[3].slice(1).toLowerCase() +
            admonitionEnd +
            check[4];

          onError({
            "lineNumber": relevantToken.lineNumber,
            "detail": relevantToken.line.substr(0, 10),
            "fixInfo": {
              "lineNumber": relevantToken.lineNumber,
              "editColumn": position,
              "deleteCount": lengthForReplacement,

              "insertText": replacement
            }
          });
        }
      };

    const relevantTokens = [];
    for (let i = 0; i < params.tokens.length; i++) {
      const token = params.tokens[i];
      if (token.type === "paragraph_open" &&
        params.tokens[i + 1].type === "inline") {
        relevantTokens.push(params.tokens[i + 1]);
      }
    }

    for (const relevantToken of relevantTokens) {
      atBeginning(relevantToken);
      wholeLine(relevantToken);
    }
  }
};

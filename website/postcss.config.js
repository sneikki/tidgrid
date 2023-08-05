module.exports = {
  parser: "postcss-scss",
  plugins: [
    require("postcss-import")({
      addModulesDirectories: ["_sass", "../build"]
    }),
    require("autoprefixer"),
    require("cssnano")({ preset: "default" })
  ]
};

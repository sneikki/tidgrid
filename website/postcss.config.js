module.exports = {
  parser: "postcss-scss",
  plugins: [
    require("postcss-import")({
      addModulesDirectories: ["_sass"]
    }),
    require("autoprefixer"),
    require("cssnano")({ preset: "default" })
  ]
};

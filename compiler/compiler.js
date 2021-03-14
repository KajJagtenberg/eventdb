var compile = function compile(src) {
  return Babel.transform(src, {
    presets: ["env"],
  }).code;
};

var compile = function compile(src) {
  return Babel.transform(src, {
    presets: ["es2016"],
  }).code;
};

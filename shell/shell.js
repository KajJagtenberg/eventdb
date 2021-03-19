Object.freeze(console);

var module = {
  exports: {},
};

var exports = {};

// const uuid = require("uuid/dist/index");

const projection = {
  ProductAdded: (state, event) => {
    state.name = event.name;
  },
};

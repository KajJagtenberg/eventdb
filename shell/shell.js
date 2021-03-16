Object.freeze(console);

const projection = {
  ProductAdded: (state, event) => {
    state.name = event.name;
  },
};

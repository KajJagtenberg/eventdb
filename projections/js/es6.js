project({
  ProductAdded: (state, event) => {
    state.name = event.name;
  },
});

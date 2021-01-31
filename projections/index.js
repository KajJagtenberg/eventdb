/***************/

const project = (projector) => {
  const state = {};
  let version = 0;

  const events = [];

  for (const event of events) {
    const handler = projector[event.type];

    if (handler) {
      handler(state, event);
    }
  }

  return {
    version,
    ...state,
  };
};

/***************/

project({
  product_added: (state, event) => {
    const { data } = event;

    state.name = data.name;
    state.price = data.price;
  },
});

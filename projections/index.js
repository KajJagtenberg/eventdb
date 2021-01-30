const events = [
  {
    stream: "3973ae68-3d49-4f8e-a03f-db2072266dbd",
    type: "product_added",
    version: 0,
    data: { name: "appel", price: 100 },
  },
  {
    stream: "3973ae68-3d49-4f8e-a03f-db2072266dbd",
    type: "product_name_changed",
    version: 1,
    data: { name: "appels" },
  },
  {
    stream: "3973ae68-3d49-4f8e-a03f-db2072266dbd",
    type: "product_price_changed",
    version: 2,
    data: { price: 150 },
  },
];

const product = events.reduce(
  (state, event) => {
    const { type, data } = event;
    switch (type) {
      case "product_added":
        state.name = data.name;
        state.price = data.price;
        break;
      case "product_name_changed":
        state.name = data.name;
        break;
      case "product_price_changed":
        state.price = data.price;
        break;
    }

    state.version++;

    return state;
  },
  {
    version: 0,
  }
);

log(JSON.stringify(product, null, 2));

when({
  product_added: (event) => {
    const { data, stream: id } = event;
    const { name, price } = data;

    set('products', id, {
      version: 1,
      name,
      price,
    });
  },
  order_added: (event) => {
    const { data, stream: id } = event;
    const { code, products, quantity } = data;

    set('orders', id, {
      code,
      products,
      quantity,
    });
  },
  $any: (event) => {
    println(event.type);
  },
});

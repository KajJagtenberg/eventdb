const projection = {
  TestEvent: (s, e) => {
    s.ids.push(e.id);
  },
  $initial: {
    ids: [],
  },
};

const project = (proj) => {
  const events = log();

  return events.reduce((s, e) => {
    const handler = proj[e.type];

    if (handler) {
      handler(s, e);
    }

    return s;
  }, projection.$initial);
};

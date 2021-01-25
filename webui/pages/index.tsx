import { useEffect, useState } from "react";
import { backend } from "../vars/backend";

const Home = () => {
  const [data, setData] = useState(null);

  useEffect(() => {
    fetch(`${backend}/api/v1`)
      .then((response) => response.json())
      .then(setData);
  }, []);

  return (
    <div>
      EventDB
      <pre>{JSON.stringify(data, null, 2)}</pre>
    </div>
  );
};

export default Home;

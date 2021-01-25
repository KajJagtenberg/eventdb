import { Box } from "@chakra-ui/react";
import { useEffect, useState } from "react";
import { useQuery } from "react-query";
import { backend } from "../vars/backend";

const fetchHome = async () => {
  const response = await fetch(`${backend}/api/v1`);
  return await response.json();
};

const Home = () => {
  const { data, isLoading } = useQuery("home", fetchHome);

  return <Box>{!isLoading && <pre>{JSON.stringify(data, null, 2)}</pre>}</Box>;
};

export default Home;

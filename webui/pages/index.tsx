import { useQuery } from 'react-query';

import { Spinner, useToast } from '@chakra-ui/react';

import Layout from '../components/Layout';
import { backend } from '../vars/backend';

const fetchHome = async () => {
  const response = await fetch(`${backend}/api/v1`);
  return response.json();
};

const Home = () => {
  const { data, isLoading } = useQuery('home', fetchHome, {
    onError: (error: Error) => {
      toast({
        title: 'Oops',
        description: error.message,
        status: 'error',
        isClosable: true,
      });
    },
  });

  const toast = useToast();

  return (
    <Layout title="EventDB">
      {data && <pre>{JSON.stringify(data, null, 2)}</pre>}
      {isLoading && <Spinner color="teal.500" />}
    </Layout>
  );
};

export default Home;

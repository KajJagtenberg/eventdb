import { Flex } from '@chakra-ui/react';
import { Layout } from '../../components/Layout';
import { Navbar } from '../../components/Navbar';
import { StreamsTable } from '../../components/StreamsTable';
import { useRouter } from 'next/router';

const Stream = () => {
  const router = useRouter();
  const { id } = router.query;

  return (
    <Layout title={`Stream - ${id}`}>
      <Flex flexDirection="column" bg="brand.500" h="100vh">
        <Navbar />

        <pre>{id}</pre>
      </Flex>
    </Layout>
  );
};

export default Stream;

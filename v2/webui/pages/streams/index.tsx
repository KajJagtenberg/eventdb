import { Flex } from '@chakra-ui/react';
import { Layout } from '../../components/Layout';
import { Navbar } from '../../components/Navbar';
import { StreamsTable } from '../../components/StreamsTable';

const Streams = () => {
  return (
    <Layout title="EventflowDB - Streams">
      <Flex flexDirection="column" bg="brand.500" h="100vh">
        <Navbar />

        <StreamsTable />
      </Flex>
    </Layout>
  );
};

export default Streams;

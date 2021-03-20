import { Flex } from '@chakra-ui/react';
import { Layout } from '../components/Layout';
import { Navbar } from '../components/Navbar';

const Dashboard = () => {
  return (
    <Layout title="EventflowDB">
      <Flex h="100vh" bg="brand.500" flexDirection="column">
        <Navbar />
      </Flex>
    </Layout>
  );
};

export default Dashboard;

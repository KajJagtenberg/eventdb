import { Layout } from '../../components/Layout';
import { Navbar } from '../../components/Navbar';
import { Flex, Heading, Text } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import { StreamTable } from '../../components/StreamTable';

const Stream = () => {
  const router = useRouter();
  const { id } = router.query;

  return (
    <Layout title={id ? `Stream - ${id}` : 'Stream'}>
      <Flex flexDirection="column" bg="brand.500" h="100vh">
        <Navbar />

        <Flex m={4} p={4} bg="white" shadow="md" flexDirection="column">
          <Flex mb={4} alignItems="center">
            <Heading size="md" mr={2} color="brand.400">
              Stream:
            </Heading>
            <Text fontWeight="semibold" color="brand.300">
              {id}
            </Text>
          </Flex>

          <StreamTable id={id} />
        </Flex>
      </Flex>
    </Layout>
  );
};

export default Stream;

import { Layout } from '../../components/Layout';
import { Navbar } from '../../components/Navbar';
import { Flex, Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import { useRouter } from 'next/router';
import { gql, useQuery } from '@apollo/client';
import Link from 'next/link';

const Stream = () => {
  const router = useRouter();
  const { id } = router.query;

  const { data } = useQuery(
    gql`
      query Get($stream: String!) {
        get(input: { stream: $stream, version: 0, limit: 0 }) {
          id
          version
          type
          data
          metadata
          causation_id
          correlation_id
          added_at
        }
      }
    `,
    {
      variables: {
        stream: id,
      },
    }
  );

  return (
    <Layout title={id ? `Stream - ${id}` : 'Stream'}>
      <Flex flexDirection="column" bg="brand.500" h="100vh">
        <Navbar />

        <Flex m={4} p={4} bg="white" shadow="md">
          <Table variant="striped" size="sm">
            <Thead>
              <Tr>
                <Th>ID</Th>
                <Th>Version</Th>
                <Th>Type</Th>
                <Th>Data</Th>
                <Th>Metadata</Th>
                <Th>Causation ID</Th>
                <Th>Correlation ID</Th>
                <Th>Added At</Th>
              </Tr>
            </Thead>
            <Tbody>
              {data?.get.map(
                (
                  {
                    id,
                    version,
                    type,
                    data,
                    metadata,
                    causation_id,
                    correlation_id,
                    added_at,
                  },
                  index
                ) => (
                  <Tr key={index}>
                    <Td>
                      <Link href={`/streams/${id}`}>
                        <a>{id}</a>
                      </Link>
                    </Td>
                    <Td>{version}</Td>
                    <Td>{type}</Td>
                    <Td>{atob(data)}</Td>
                    <Td>{atob(metadata) || '-'}</Td>
                    <Td>{causation_id}</Td>
                    <Td>{correlation_id}</Td>
                    <Td>{added_at}</Td>
                  </Tr>
                )
              )}
            </Tbody>
          </Table>
        </Flex>
      </Flex>
    </Layout>
  );
};

export default Stream;

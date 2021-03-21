import { Layout } from '../../components/Layout';
import { Navbar } from '../../components/Navbar';
import {
  Flex,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Link as UILink,
  Heading,
  Text,
} from '@chakra-ui/react';
import { useRouter } from 'next/router';
import { gql, useQuery } from '@apollo/client';
import Link from 'next/link';
import dayjs from 'dayjs';

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

        <Flex m={4} p={4} bg="white" shadow="md" flexDirection="column">
          <Flex mb={4} alignItems="center">
            <Heading size="md" mr={2} color="brand.400">
              Stream:
            </Heading>
            <Text fontWeight="semibold" color="brand.300">{id}</Text>
          </Flex>

          <Flex>
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
                    }: {
                      id: string;
                      version: number;
                      type: string;
                      data: string;
                      metadata: string;
                      causation_id: string;
                      correlation_id: string;
                      added_at: number;
                    },
                    index: number
                  ) => (
                    <Tr key={index}>
                      <Td>
                        <UILink color="blue.400">
                          <Link href={`/events/${id}`}>
                            <a>{id}</a>
                          </Link>
                        </UILink>
                      </Td>
                      <Td>{version}</Td>
                      <Td>{type}</Td>
                      <Td>{atob(data)}</Td>
                      <Td>{atob(metadata) || '-'}</Td>
                      <Td>
                        <UILink color="blue.400">
                          <Link href={`/events/${causation_id}`}>
                            <a>{causation_id}</a>
                          </Link>
                        </UILink>
                      </Td>
                      <Td>
                        <UILink color="blue.400">
                          <Link href={`/events/${correlation_id}`}>
                            <a>{correlation_id}</a>
                          </Link>
                        </UILink>
                      </Td>
                      <Td>{dayjs(added_at).format('HH:mm:ss DD-MM-YYYY')}</Td>
                    </Tr>
                  )
                )}
              </Tbody>
            </Table>
          </Flex>
        </Flex>
      </Flex>
    </Layout>
  );
};

export default Stream;

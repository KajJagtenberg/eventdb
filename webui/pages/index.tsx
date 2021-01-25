import Link from 'next/link';
import { useState } from 'react';
import { useQuery } from 'react-query';

import {
    Box, Button, Flex, Link as UILink, Spinner, Table, TableCaption, Tbody, Td, Text, Tfoot, Th,
    Thead, Tr
} from '@chakra-ui/react';

import Layout from '../components/Layout';
import { backend } from '../vars/backend';

const fetchStreams = async (page: number, limit: number) => {
  const response = await fetch(
    `${backend}/api/v1/streams?offset=${(page - 1) * limit}&limit=${limit}`
  );
  return response.json();
};

const Home = () => {
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const { data } = useQuery(
    ['streams', page, limit],
    () => fetchStreams(page, limit),
    { keepPreviousData: true }
  );

  const lastPage = data ? Math.floor(data.total / limit) + 1 : 1;

  return (
    <Layout title="EventDB">
      <Flex
        direction="column"
        p={4}
        m={4}
        bg="white"
        style={{
          boxShadow: '0 0 10px #aaa',
        }}
      >
        <Flex direction="row" alignItems="center" mb={8}>
          <Button
            colorScheme="teal"
            size="sm"
            mx={1}
            minW={20}
            onClick={() => setPage(1)}
            disabled={page === 1}
          >
            First
          </Button>

          <Button
            colorScheme="teal"
            size="sm"
            mx={1}
            minW={20}
            onClick={() => setPage((old) => Math.max(1, old - 1))}
            disabled={page === 1}
          >
            Previous
          </Button>

          <Text fontWeight="semibold" mx={1}>
            {page}/{lastPage}
          </Text>

          <Button
            colorScheme="teal"
            size="sm"
            mx={1}
            minW={20}
            onClick={() => setPage((old) => old + 1)}
            disabled={page === lastPage}
          >
            Next
          </Button>

          <Button
            colorScheme="teal"
            size="sm"
            mx={1}
            minW={20}
            onClick={() => setPage(lastPage)}
            disabled={page === lastPage}
          >
            Last
          </Button>
        </Flex>

        <Table variant="striped" size="sm">
          <Thead>
            <Tr>
              <Th>Row</Th>
              <Th>Stream</Th>
            </Tr>
          </Thead>

          <Tbody>
            {data &&
              data.streams.map((stream: string, index: number) => {
                return (
                  <Tr key={index}>
                    <Td>{index + 1}</Td>
                    <Td>
                      <Link href={`${backend}/streams/${stream}`}>
                        <UILink textDecoration="underline" color="blue.600">
                          {stream}
                        </UILink>
                      </Link>
                    </Td>
                  </Tr>
                );
              })}
            {!data && <Spinner color="teal.500" />}
          </Tbody>
        </Table>
      </Flex>
    </Layout>
  );
};

export default Home;

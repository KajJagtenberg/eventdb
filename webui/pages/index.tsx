import Link from 'next/link';
import { useState } from 'react';
import { useQuery } from 'react-query';

import {
    Box, Button, Flex, Link as UILink, Spinner, Table, TableCaption, Tbody, Td, Text, Tfoot, Th,
    Thead, Tr, useToast
} from '@chakra-ui/react';

import Layout from '../components/Layout';
import Navbar from '../components/Navbar';
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
  const { data, isLoading } = useQuery(
    ['streams', page, limit],
    () => fetchStreams(page, limit),
    {
      keepPreviousData: true,
      onError: (error: Error) => {
        toast({
          title: 'Oops',
          description: error.message,
          status: 'error',
        });
      },
    }
  );

  const lastPage = data ? Math.floor(data.total / limit) + 1 : 1;

  const toast = useToast();

  return (
    <Layout title="EventDB">
      <Navbar />

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
              <Th>#</Th>
              <Th>Stream</Th>
            </Tr>
          </Thead>

          <Tbody>
            {data &&
              data.streams.map((stream: string, index: number) => {
                return (
                  <Tr key={index}>
                    <Td>{(page - 1) * limit + index + 1}</Td>
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
          </Tbody>
        </Table>

        {isLoading && (
          <Flex w="full" justifyContent="center" my={8} color="teal.500">
            <Spinner size="lg" />
          </Flex>
        )}
      </Flex>
    </Layout>
  );
};

export default Home;

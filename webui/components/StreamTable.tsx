import Link from 'next/link';
import { useState } from 'react';

import {
  Box,
  Button,
  Flex,
  Link as UILink,
  Spinner,
  Table,
  Tbody,
  Td,
  Text,
  Th,
  Thead,
  Tr,
  useToast,
} from '@chakra-ui/react';

import { useQuery, gql } from '@apollo/client';

const StreamTable = () => {
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const { loading, error, data } = useQuery(
    gql`
      query Streams($skip: Int!, $limit: Int!) {
        streams(skip: $skip, limit: $limit) {
          name
        }
      }
    `,
    {
      variables: {
        skip: (page - 1) * limit * 0,
        limit: 0,
      },
    }
  );

  const lastPage = data ? Math.floor(data.total / limit) + 1 : 1;

  const toast = useToast();

  if (error) {
    toast({
      title: 'Oops',
      description: error.message,
      status: 'error',
    });
  }

  return (
    <Flex
      direction="column"
      p={4}
      m={4}
      bg="white"
      style={{
        boxShadow: '0 0 10px #aaa',
      }}
    >
      <Flex
        direction="row"
        alignItems="center"
        justifyContent="space-between"
        mb={8}
      >
        <Flex direction="row">
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

          <pre>{JSON.stringify(loading)}</pre>

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

        <Flex>
          {/* <Button onClick={onOpen} size="sm" colorScheme="teal">
            Add Event
          </Button> */}
        </Flex>
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
            data.streams.map(({ name }: { name: string }, index: number) => {
              return (
                <Tr key={index}>
                  <Td>{(page - 1) * limit + index + 1}</Td>
                  <Td>
                    <UILink textDecoration="underline" color="blue.600">
                      <Link href={`/streams/${name}`}>
                        <a>{name}</a>
                      </Link>
                    </UILink>
                  </Td>
                </Tr>
              );
            })}
        </Tbody>
      </Table>

      {data && data.streams.length === 0 && (
        <Box textAlign="center" w="full" my={8}>
          <Text fontWeight="semibold">No streams available</Text>
        </Box>
      )}

      {loading && (
        <Flex w="full" justifyContent="center" my={8} color="teal.500">
          <Spinner size="lg" />
        </Flex>
      )}

      {/* <AddEventModal isOpen={isOpen} onClose={onClose} /> */}
    </Flex>
  );
};

export default StreamTable;

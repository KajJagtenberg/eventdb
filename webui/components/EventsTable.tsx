import { useState } from 'react';
import { useQuery } from 'react-query';

import {
  Button,
  Flex,
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

import { backend } from '../vars/backend';
import dayjs from 'dayjs';

const fetchStreams = async (stream: string, page: number, limit: number) => {
  const response = await fetch(
    `${backend}/api/v1/streams/${stream}?offset=${
      (page - 1) * limit
    }&limit=${limit}`
  );
  return response.json();
};

const EventsTable = ({ stream }) => {
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const { data, isLoading } = useQuery(
    ['events', page, limit],
    () => fetchStreams(stream as string, page, limit),
    {
      // keepPreviousData: true,
      onError: (error: Error) => {
        toast({
          title: 'Oops',
          description: error.message,
          isClosable: true,
          status: 'error',
        });
      },
    }
  );

  const lastPage = data ? Math.floor(data.total / limit) + 1 : 1;

  const toast = useToast();
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
            <Th>Version</Th>
            <Th>ID</Th>
            <Th>Type</Th>
            <Th>Timestamp</Th>
            <Th>Payload</Th>
          </Tr>
        </Thead>

        <Tbody>
          {!isLoading &&
            data &&
            data.events.map(
              ({ version, id, type, ts, data }, index: number) => {
                return (
                  <Tr key={index}>
                    <Td>{version}</Td>
                    <Td>{id}</Td>
                    <Td>{type}</Td>
                    <Td>{dayjs(ts).format('DD-MM-YYYY HH:mm:ss:ms')}</Td>
                    <Td>{JSON.stringify(data)}</Td>
                  </Tr>
                );
              }
            )}
        </Tbody>
      </Table>

      {isLoading && (
        <Flex w="full" justifyContent="center" my={8} color="teal.500">
          <Spinner size="lg" />
        </Flex>
      )}
    </Flex>
  );
};

export default EventsTable;

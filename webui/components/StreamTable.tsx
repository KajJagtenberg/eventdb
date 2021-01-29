import Link from 'next/link';
import { useState } from 'react';
import { useQuery } from 'react-query';

import {
    Box, Button, Flex, Link as UILink, Spinner, Table, Tbody, Td, Text, Th, Thead, Tr,
    useDisclosure, useToast
} from '@chakra-ui/react';

import { backend } from '../vars/backend';
import AddEventModal from './AddEventModal';

const fetchStreams = async (page: number, limit: number) => {
  const response = await fetch(
    `${backend}/api/v1/streams?offset=${(page - 1) * limit}&limit=${limit}`
  );
  return response.json();
};

const StreamTable = () => {
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
          isClosable: true,
          status: 'error',
        });
      },
    }
  );

  const lastPage = data ? Math.floor(data.total / limit) + 1 : 1;

  const toast = useToast();

  const { isOpen, onOpen, onClose } = useDisclosure();

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
            data.streams.map((stream: string, index: number) => {
              return (
                <Tr key={index}>
                  <Td>{(page - 1) * limit + index + 1}</Td>
                  <Td>
                    <UILink textDecoration="underline" color="blue.600">
                      <Link href={`/streams/${stream}`}>
                        <a>{stream}</a>
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

      {isLoading && (
        <Flex w="full" justifyContent="center" my={8} color="teal.500">
          <Spinner size="lg" />
        </Flex>
      )}

      <AddEventModal isOpen={isOpen} onClose={onClose} />
    </Flex>
  );
};

export default StreamTable;

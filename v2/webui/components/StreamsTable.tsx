import { gql, useQuery } from '@apollo/client';
import {
  Flex,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Link as UILink,
  Button,
  Text,
} from '@chakra-ui/react';
import Link from 'next/link';
import dayjs from 'dayjs';
import { useState } from 'react';

export const StreamsTable = () => {
  const [page, setPage] = useState(0);
  const [limit, setLimit] = useState(20);

  const { data, previousData } = useQuery(
    gql`
      query Streams($skip: Int! = 0, $limit: Int! = 20) {
        streams(input: { skip: $skip, limit: $limit }) {
          id
          added_at
        }
        streamCount
      }
    `,
    {
      variables: {
        skip: page * limit,
        limit,
      },
      pollInterval: 1000,
    }
  );

  const result = data || previousData;

  const lastPage = Math.floor(result?.streamCount / limit);

  return (
    <Flex m={4} p={4} bg="white" flexDirection="column">
      <Flex mb={2} alignItems="center">
        <Button
          size="sm"
          colorScheme="brand"
          mx={1}
          disabled={page == 0}
          onClick={() => {
            setPage(0);
          }}
        >
          First
        </Button>
        <Button
          size="sm"
          colorScheme="brand"
          mx={1}
          disabled={page == 0}
          onClick={() => {
            setPage((old) => {
              return Math.max(old - 1, 0);
            });
          }}
        >
          Previous
        </Button>

        <Text fontWeight="semibold" mx={1}>
          {page + 1} / {lastPage + 1}
        </Text>

        <Button
          size="sm"
          colorScheme="brand"
          mx={1}
          disabled={page == lastPage}
          onClick={() => {
            setPage((old) => {
              return Math.min(old + 1, lastPage);
            });
          }}
        >
          Next
        </Button>
        <Button
          size="sm"
          colorScheme="brand"
          mx={1}
          disabled={page == lastPage}
          onClick={() => {
            setPage(lastPage);
          }}
        >
          Last
        </Button>
      </Flex>

      <Flex>
        <Table variant="simple" size="sm">
          <Thead>
            <Tr>
              <Th>ID</Th>
              <Th>Added At (D/M/Y)</Th>
            </Tr>
          </Thead>
          <Tbody>
            {result?.streams.map(({ id, added_at }, index) => (
              <Tr key={index}>
                <Td>
                  <UILink color="blue.400">
                    <Link href={`/streams/${id}`}>
                      <a>{id}</a>
                    </Link>
                  </UILink>
                </Td>
                <Td>{dayjs(added_at).format('HH:mm:ss DD-MM-YYYY')}</Td>
              </Tr>
            ))}
          </Tbody>
        </Table>
      </Flex>
    </Flex>
  );
};

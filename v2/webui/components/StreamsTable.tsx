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
} from '@chakra-ui/react';
import Link from 'next/link';
import dayjs from 'dayjs';

export const StreamsTable = () => {
  const { data } = useQuery(gql`
    {
      streams {
        id
        added_at
      }
    }
  `);

  return (
    <Flex m={4} p={4} bg="white">
      <Table variant="simple" size="sm">
        <Thead>
          <Tr>
            <Th>ID</Th>
            <Th>Added At (D/M/Y)</Th>
          </Tr>
        </Thead>
        <Tbody>
          {data?.streams.map(({ id, added_at }, index) => (
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
  );
};

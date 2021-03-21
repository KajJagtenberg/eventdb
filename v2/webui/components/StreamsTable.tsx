import { gql, useQuery } from '@apollo/client';
import { Flex, Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';
import Link from 'next/link';

export const StreamsTable = () => {
  const { data } = useQuery(gql`
    {
      streams {
        id
      }
    }
  `);

  return (
    <Flex m={4} p={4} bg="white">
      <Table variant="simple">
        <Thead size="sm">
          <Tr>
            <Th>ID</Th>
            <Th>Size</Th>
          </Tr>
        </Thead>
        <Tbody>
          {data?.streams.map(({ id }, index) => (
            <Tr key={index}>
              <Td>
                <Link href={`/streams/${id}`}>
                  <a>{id}</a>
                </Link>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Flex>
  );
};

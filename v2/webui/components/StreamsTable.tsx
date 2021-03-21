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
      <Table variant="simple" size="sm">
        <Thead>
          <Tr>
            <Th>ID</Th>
            {/* <Th>Size</Th> */}
          </Tr>
        </Thead>
        <Tbody>
          {data?.streams.map(({ id }, index) => (
            <Tr key={index}>
              <Td>
                <UILink color="blue.400">
                  <Link href={`/streams/${id}`}>
                    <a>{id}</a>
                  </Link>
                </UILink>
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Flex>
  );
};

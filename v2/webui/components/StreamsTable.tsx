import { gql, useQuery } from '@apollo/client';
import { Flex, Table, Thead, Tbody, Tr, Th, Td } from '@chakra-ui/react';

export const StreamsTable = () => {
  const { data } = useQuery(gql`
    {
      nodes {
        ip
        port
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
          {data?.nodes.map(({ ip, port }, index) => (
            <Tr key={index}>
              <Td>{ip}</Td>
              <Td>{port}</Td>
            </Tr>
          ))}
        </Tbody>
      </Table>
    </Flex>
  );
};

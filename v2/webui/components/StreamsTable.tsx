import { Flex, Table, Thead, Tbody, Tr, Th } from '@chakra-ui/react';

export const StreamsTable = () => {
  return (
    <Flex m={4} p={4} bg="white">
      <Table variant="simple">
        <Thead size="sm">
          <Tr>
            <Th>ID</Th>
            <Th>Size</Th>
          </Tr>
        </Thead>
        <Tbody></Tbody>
      </Table>
    </Flex>
  );
};

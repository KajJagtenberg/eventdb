import { gql, useQuery } from '@apollo/client';
import { Flex, Text } from '@chakra-ui/react';

export const ClusterTab = () => {
  const { data } = useQuery(gql`
    {
      nodes {
        address
      }
    }
  `);

  return (
    <>
      <Flex>
        <Flex>
          <Text fontWeight="semibold" mr={2} color="brand.500">
            Nodes:
          </Text>
        </Flex>

        <Flex flexDirection="column">
          {data?.nodes.map(({ address }) => (
            <Text>{address}</Text>
          ))}
        </Flex>
      </Flex>
    </>
  );
};

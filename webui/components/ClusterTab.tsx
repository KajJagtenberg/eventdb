import { gql, useQuery } from '@apollo/client';
import { Flex, Text } from '@chakra-ui/react';

export const ClusterTab = () => {
  const { data } = useQuery(
    gql`
      {
        cluster {
          leader
          size
        }
      }
    `,
    {
      pollInterval: 1000,
    }
  );

  return (
    <>
      <Flex>
        <Text fontWeight="semibold" mr={2} color="brand.500">
          Leader:
        </Text>
        <Text>{data?.cluster.leader}</Text>
      </Flex>

      <Flex>
        <Text fontWeight="semibold" mr={2} color="brand.500">
          Size:
        </Text>
        <Text>{data?.cluster.size}</Text>
      </Flex>
    </>
  );
};

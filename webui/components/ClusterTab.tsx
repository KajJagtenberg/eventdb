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
        <Flex>
          <Text fontWeight="semibold" mr={2} color="brand.500">
            Leader:
          </Text>
        </Flex>

        <Flex flexDirection="column">{data?.cluster.leader}</Flex>
      </Flex>

      <Flex>
        <Flex>
          <Text fontWeight="semibold" mr={2} color="brand.500">
            Size:
          </Text>
        </Flex>

        <Flex flexDirection="column">{data?.cluster.size}</Flex>
      </Flex>
    </>
  );
};

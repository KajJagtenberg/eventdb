import { gql, useQuery } from '@apollo/client';
import { Flex, TabPanel, Text } from '@chakra-ui/react';

export const OverviewTab = () => {
  const { data } = useQuery(
    gql`
      {
        diagnostics {
          uptime
          uptimeMs
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
          Uptime:
        </Text>
        <Text>{data?.diagnostics.uptime}(s)</Text>
      </Flex>

      <Flex>
        <Text fontWeight="semibold" mr={2} color="brand.500">
          Event Count:
        </Text>
        <Text>{data?.diagnostics.eventCount}</Text>
      </Flex>
    </>
  );
};

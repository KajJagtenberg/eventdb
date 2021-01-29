import { Field, Form, Formik } from 'formik';

import {
    Box, Button, Flex, Input, Modal, ModalBody, ModalContent, ModalFooter, ModalHeader,
    ModalOverlay, Text, Textarea
} from '@chakra-ui/react';

const AddEventModal = ({ isOpen, onClose }) => {
  return (
    <Modal isOpen={isOpen} onClose={onClose} size="2xl">
      <ModalOverlay />
      <ModalContent>
        <Formik
          initialValues={{
            stream: '',
            version: 0,
            type: '',
            data: null,
            causation_id: '',
            correlation_id: '',
          }}
          onSubmit={() => {}}
        >
          {({ values }) => {
            return (
              <Form>
                <ModalHeader>Add event</ModalHeader>
                <ModalBody>
                  <Flex direction="row" my={4}>
                    <Box mr={2}>
                      <Text fontWeight="semibold" mb={2}>
                        Stream
                      </Text>
                      <Field
                        name="stream"
                        as={Input}
                        placeholder="e.g. 58999445-0d45-4fa4-96c8-0f54c0e6e905"
                      />
                    </Box>

                    <Box>
                      <Text fontWeight="semibold" mb={2}>
                        Version
                      </Text>
                      <Field
                        name="version"
                        type="number"
                        as={Input}
                        placeholder="e.g. 58999445-0d45-4fa4-96c8-0f54c0e6e905"
                      />
                    </Box>
                  </Flex>

                  <Flex direction="row" my={4}>
                    <Box mx={1}>
                      <Text fontWeight="semibold" mb={2}>
                        Type
                      </Text>
                      <Field
                        name="stream"
                        as={Input}
                        placeholder="e.g. product_added"
                      />
                    </Box>

                    <Box mx={1}>
                      <Text fontWeight="semibold" mb={2}>
                        Causation ID
                      </Text>
                      <Field
                        name="causation_id"
                        type="text"
                        as={Input}
                        placeholder="e.g. bc66c483-899f-4fe5-8f27-32b8355ebc20"
                      />
                    </Box>

                    <Box mx={1}>
                      <Text fontWeight="semibold" mb={2}>
                        Correlation ID
                      </Text>
                      <Field
                        name="correlation_id"
                        type="text"
                        as={Input}
                        placeholder="e.g. ba17800a-ab59-4ca6-a7d4-d4884ad62f56"
                      />
                    </Box>
                  </Flex>

                  <Flex my={2}>
                    <Box w="full">
                      <Text fontWeight="semibold" mb={2}>
                        Data
                      </Text>
                      <Field
                        name="data"
                        type="text"
                        as={Textarea}
                        resize="none"
                        h={32}
                        placeholder="e.g. {}"
                      />
                    </Box>
                  </Flex>
                </ModalBody>

                <ModalFooter>
                  <Button mr={3} onClick={onClose} size="sm">
                    Cancel
                  </Button>

                  <Button colorScheme="green" size="sm">
                    Add
                  </Button>
                </ModalFooter>
              </Form>
            );
          }}
        </Formik>
      </ModalContent>
    </Modal>
  );
};

export default AddEventModal;
